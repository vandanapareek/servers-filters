"""Script for Excel to SQLite conversion"""

import argparse
import os
import re
import sqlite3
import sys
import tempfile
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple

import pandas as pd


def parse_ram(ram_str: str) -> Optional[int]:
    """Parse RAM string to extract GB value."""
    if not ram_str or pd.isna(ram_str):
        return None

    # Extract numeric value
    match = re.search(r'(\d+)\s*GB', str(ram_str), re.IGNORECASE)
    if match:
        return int(match.group(1))
    return None


def parse_storage(hdd_str: str) -> Tuple[Optional[int], str]:
    """Parse HDD string"""
    if not hdd_str or pd.isna(hdd_str):
        return None, str(hdd_str) if hdd_str else ""

    hdd_str = str(hdd_str)

    # Pattern to match: NxXTB or NxXGB
    pattern = r'(\d+)x(\d+)(TB|GB)'
    matches = re.findall(pattern, hdd_str, re.IGNORECASE)

    if matches:
        total_gb = 0
        for count, size, unit in matches:
            size_gb = int(size)
            if unit.upper() == 'TB':
                size_gb *= 1024
            total_gb += int(count) * size_gb
        return total_gb, hdd_str

    # Try single value pattern
    single_match = re.search(r'(\d+)(TB|GB)', hdd_str, re.IGNORECASE)
    if single_match:
        size = int(single_match.group(1))
        unit = single_match.group(2).upper()
        if unit == 'TB':
            size *= 1024
        return size, hdd_str

    return None, hdd_str


def parse_cpu(model_str: str) -> Optional[str]:
    """Extract CPU information"""
    if not model_str or pd.isna(model_str):
        return None

    model_str = str(model_str)

    # Common CPU patterns
    cpu_patterns = [
        r'(Intel\s+\w+)',
        r'(AMD\s+\w+)',
        r'(Xeon\s+\w+)',
        r'(Core\s+i\d+)',
        r'(Ryzen\s+\w+)',
    ]

    for pattern in cpu_patterns:
        match = re.search(pattern, model_str, re.IGNORECASE)
        if match:
            return match.group(1).strip()

    return None


def parse_price(price_str: str) -> Tuple[Optional[float], str]:
    """Parse price string"""
    if not price_str or pd.isna(price_str):
        return None, str(price_str) if price_str else ""

    price_str = str(price_str)

    # Remove currency symbols and extract numeric value ("€123.45", "123.45€", "123,45€").
    cleaned = re.sub(r'[€$£¥,\s]', '', price_str)

    # Extract decimal number
    match = re.search(r'(\d+\.?\d*)', cleaned)
    if match:
        try:
            return float(match.group(1)), price_str
        except ValueError:
            pass

    return None, price_str


def parse_location(location_str: str) -> Tuple[Optional[str], Optional[str]]:
    """Parse location string"""
    if not location_str or pd.isna(location_str):
        return None, None

    location_str = str(location_str)

    # Pattern for city-code format - updated to handle periods and 3-letter codes
    # Examples: "AmsterdamAMS-01", "Washington D.C.WDC-01", "FrankfurtFRA-10"
    match = re.match(r'^([A-Za-z\s\.]+?)([A-Z]{2,4}-\d+)$', location_str)
    if match:
        city = match.group(1).strip()
        code = match.group(2).strip()
        return city, code

    # If no code pattern, treat as city
    return location_str.strip(), None


def create_database_schema(conn: sqlite3.Connection) -> None:
    """Create the servers table"""
    cursor = conn.cursor()

    # Drop table if exists
    cursor.execute("DROP TABLE IF EXISTS servers")

    # Create servers table
    cursor.execute("""
        CREATE TABLE servers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            model TEXT NOT NULL,
            cpu TEXT,
            ram_gb INTEGER,
            hdd TEXT,
            storage_gb INTEGER,
            location_city TEXT,
            location_code TEXT,
            price_eur REAL,
            raw_price TEXT,
            raw_ram TEXT,
            raw_hdd TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)

    # Create indexes
    indexes = [
        "CREATE INDEX idx_servers_model ON servers(model)",
        "CREATE INDEX idx_servers_cpu ON servers(cpu)",
        "CREATE INDEX idx_servers_ram_gb ON servers(ram_gb)",
        "CREATE INDEX idx_servers_storage_gb ON servers(storage_gb)",
        "CREATE INDEX idx_servers_location_city ON servers(location_city)",
        "CREATE INDEX idx_servers_location_code ON servers(location_code)",
        "CREATE INDEX idx_servers_price_eur ON servers(price_eur)",
        "CREATE INDEX idx_servers_hdd ON servers(hdd)",
    ]

    for index_sql in indexes:
        cursor.execute(index_sql)

    conn.commit()


def process_excel_data(df: pd.DataFrame) -> List[Dict[str, Any]]:
    """Process Excel data"""
    processed_data = []

    for _, row in df.iterrows():
        # Parse RAM
        ram_gb = parse_ram(row.get('RAM', ''))

        # Parse storage
        storage_gb, raw_hdd = parse_storage(row.get('HDD', ''))

        # Parse CPU from model
        cpu = parse_cpu(row.get('Model', ''))

        # Parse price
        price_eur, raw_price = parse_price(row.get('Price', ''))

        # Parse location
        location_city, location_code = parse_location(row.get('Location', ''))

        processed_row = {
            'model': str(row.get('Model', '')),
            'cpu': cpu,
            'ram_gb': ram_gb,
            'hdd': str(row.get('HDD', '')),
            'storage_gb': storage_gb,
            'location_city': location_city,
            'location_code': location_code,
            'price_eur': price_eur,
            'raw_price': raw_price,
            'raw_ram': str(row.get('RAM', '')),
            'raw_hdd': raw_hdd,
        }

        processed_data.append(processed_row)

    return processed_data


def insert_data_batch(conn: sqlite3.Connection, data: List[Dict[str, Any]]) -> None:
    """Insert processed data into the database in batches."""
    cursor = conn.cursor()

    # Prepare insert statement
    insert_sql = """
        INSERT INTO servers (
            model, cpu, ram_gb, hdd, storage_gb, location_city, 
            location_code, price_eur, raw_price, raw_ram, raw_hdd
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """

    # Convert data to tuples for batch insert
    data_tuples = [
        (
            row['model'], row['cpu'], row['ram_gb'], row['hdd'],
            row['storage_gb'], row['location_city'], row['location_code'],
            row['price_eur'], row['raw_price'], row['raw_ram'], row['raw_hdd']
        )
        for row in data
    ]

    # Insert in batches for better performance
    batch_size = 1000
    for i in range(0, len(data_tuples), batch_size):
        batch = data_tuples[i:i + batch_size]
        cursor.executemany(insert_sql, batch)
        conn.commit()
        print(
            f"Inserted batch {i//batch_size + 1}/{(len(data_tuples) + batch_size - 1)//batch_size}")


def convert_excel_to_sqlite(excel_path: str, output_path: str) -> None:
    """Convert Excel file to SQLite"""
    print(f"Reading Excel file: {excel_path}")

    # Read Excel file in chunks to handle large files
    try:
        # Try to read the first sheet
        df = pd.read_excel(excel_path, engine='openpyxl')
        print(f"Loaded {len(df)} rows from Excel file")
    except Exception as e:
        print(f"Error reading Excel file: {e}")
        sys.exit(1)

    # Process the data
    print("Processing and normalizing data...")
    processed_data = process_excel_data(df)
    print(f"Processed {len(processed_data)} rows")

    # Create temporary database file
    temp_dir = Path(output_path).parent
    temp_db = tempfile.NamedTemporaryFile(
        suffix='.db',
        dir=temp_dir,
        delete=False
    )
    temp_db_path = temp_db.name
    temp_db.close()

    try:
        # Create database schema
        print("Creating database schema...")
        conn = sqlite3.connect(temp_db_path)
        create_database_schema(conn)

        # Insert data
        print("Inserting data into database...")
        insert_data_batch(conn, processed_data)

        # Verify data
        cursor = conn.cursor()
        cursor.execute("SELECT COUNT(*) FROM servers")
        count = cursor.fetchone()[0]
        print(f"Successfully inserted {count} records")

        conn.close()

        # Atomically move temp file to final location
        if os.path.exists(output_path):
            os.remove(output_path)
        os.rename(temp_db_path, output_path)

        print(f"Database created successfully: {output_path}")

    except Exception as e:
        # Clean up temp file on error
        if os.path.exists(temp_db_path):
            os.remove(temp_db_path)
        print(f"Error creating database: {e}")
        sys.exit(1)


def main():
    parser = argparse.ArgumentParser(
        description="Convert Excel server data to SQLite database"
    )
    parser.add_argument(
        "excel_file",
        help="Path to the Excel file containing server data"
    )
    parser.add_argument(
        "-o", "--output",
        default="data/servers.db",
        help="Output SQLite database path (default: data/servers.db)"
    )

    args = parser.parse_args()

    # Validate input file
    if not os.path.exists(args.excel_file):
        print(f"Error: Excel file not found: {args.excel_file}")
        sys.exit(1)

    # Create output directory if it doesn't exist
    output_dir = Path(args.output).parent
    output_dir.mkdir(parents=True, exist_ok=True)

    # Convert Excel to SQLite
    convert_excel_to_sqlite(args.excel_file, args.output)


if __name__ == "__main__":
    main()
