import pandas as pd
input_csv = r"C:\Users\aikaw\OneDrive\Documents\Dev\GeoGO\data\Datasets\Meteorite_Landings.csv"
output_csv = r"C:\Datasets\Meteorite_Landings_CLEANED.csv"
# Load CSV
df = pd.read_csv(input_csv)
# Drop unnamed columns automatically
df = df.loc[:, ~df.columns.str.contains('^Unnamed')]
# Drop "GeoLocation" column if it exists
if "GeoLocation" in df.columns:
    df = df.drop(columns=["GeoLocation"])
# Rename "mass (g)" to "mass" if it exists
if "mass (g)" in df.columns:
    df.rename(columns={"mass (g)": "mass"}, inplace=True)
# Ensure year is stored as an INTEGER
df["year"] = pd.Series(pd.to_numeric(df["year"], errors='coerce')).fillna(0).astype(int)
# **REORDER COLUMNS TO MATCH PostgreSQL TABLE**
df = df[['id', 'name', 'nametype', 'recclass', 'mass', 'fall', 'year', 'reclat', 'reclong']]
# Fill missing values with NULL
df = df.fillna("NULL")
# Save cleaned CSV
df.to_csv(output_csv, index=False)
print(f"✅ Cleaned CSV saved at: {output_csv}")
