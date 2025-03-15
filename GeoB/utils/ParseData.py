import pandas as pd
import os
DATASET_PATH = r"C:\Users\aikaw\OneDrive\Documents\Dev\GeoGO\data\Datasets\Meteorite_Landings.csv"
OUTPUT_SQL_PATH = r"C:\Users\aikaw\OneDrive\Documents\Dev\GeoGO\utils\SQL\seed_meteorites.sql"
os.makedirs(os.path.dirname(OUTPUT_SQL_PATH), exist_ok=True)
df = pd.read_csv(DATASET_PATH)
# Coluns
df = df[['name', 'id', 'nametype', 'recclass', 'mass (g)', 'fall', 'year', 'reclat', 'reclong']]
# Drop rows with missing coordinates
df = df.dropna(subset=['reclat', 'reclong'])
# Convert 'mass (g)' to float and fill NaNs with 0
df['mass (g)'] = pd.to_numeric(df['mass (g)'], errors='coerce').fillna(0)
# Format SQL insert statements
sql_statements = "INSERT INTO locations (id, name, nametype, recclass, mass, fall, year, latitude, longitude) VALUES\n"
values = []
for _, row in df.iterrows():
    id_val = int(row['id'])
    name = row['name'].replace("'", "''")  
    nametype = row['nametype'].replace("'", "''")
    recclass = row['recclass'].replace("'", "''")
    mass = row['mass (g)']
    fall = row['fall'].replace("'", "''")
    year = int(row['year']) if pd.notna(row['year']) else "NULL"
    lat = row['reclat']
    lon = row['reclong']
    values.append(f"({id_val}, '{name}', '{nametype}', '{recclass}', {mass}, '{fall}', {year}, {lat}, {lon})")
# Join values into a single SQL statement
sql_statements += ",\n".join(values) + ";\n"
# Write to SQL file
with open(OUTPUT_SQL_PATH, "w", encoding="utf-8") as f:
    f.write(sql_statements)
print(f"✅ SQL file generated: {OUTPUT_SQL_PATH}")
