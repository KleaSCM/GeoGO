import pandas as pd
import json
import os
from datetime import datetime
import sqlite3

class DatasetProcessor:
    def __init__(self, data_dir="data/Datasets", output_dir="utils/SQL"):
        self.data_dir = data_dir
        self.output_dir = output_dir
        os.makedirs(output_dir, exist_ok=True)
        
    def process_all_datasets(self):
        """Process all datasets and generate SQL files"""
        print("🚀 Starting multi-dataset processing...")
        
        # Process each dataset type
        self.process_meteorites()
        self.process_climate_data()
        self.process_wind_data()
        self.process_vegetation_data()
        
        print("✅ All datasets processed!")
    
    def process_meteorites(self):
        """Process meteorite data (existing functionality)"""
        print("☄️ Processing meteorite data...")
        input_file = os.path.join(self.data_dir, "Meteorite_Landings.csv")
        
        if not os.path.exists(input_file):
            print(f"⚠️ Skipping meteorites - file not found: {input_file}")
            return
            
        df = pd.read_csv(input_file)
        df = df.loc[:, ~df.columns.str.contains('^Unnamed')]
        
        if "GeoLocation" in df.columns:
            df = df.drop(columns=["GeoLocation"])
        if "mass (g)" in df.columns:
            df.rename(columns={"mass (g)": "mass"}, inplace=True)
            
        df["year"] = pd.Series(pd.to_numeric(df["year"], errors='coerce')).fillna(0).astype(int)
        df = df[['id', 'name', 'nametype', 'recclass', 'mass', 'fall', 'year', 'reclat', 'reclong']]
        df = df.fillna("NULL")
        
        # Convert to unified format
        unified_data = []
        for _, row in df.iterrows():
            # Skip records with missing coordinates
            if pd.isna(row['reclat']) or pd.isna(row['reclong']):
                continue
                
            metadata = {
                "recclass": row['recclass'],
                "mass": row['mass'],
                "year": row['year'],
                "nametype": row['nametype'],
                "fall": row['fall']
            }
            
            unified_data.append({
                "dataset_type": "meteorite",
                "name": row['name'],
                "lat": row['reclat'],
                "lon": row['reclong'],
                "value": row['mass'],
                "unit": "g",
                "metadata": json.dumps(metadata)
            })
        
        self.save_to_sql(unified_data, "meteorites.sql")
    
    def process_climate_data(self):
        """Process climate data files"""
        print("🌡️ Processing climate data...")
        
        climate_files = [
            "tas_aus-station_r1i1p1_CSIRO-MnCh-wrt-1986-2005-Scl_v1_mon_seasavg-clim_1.csv",
            "tasmax_aus-station_r1i1p1_CSIRO-MnCh-wrt-1986-2005-Scl_v1_mon_seasavg-clim.csv",
            "tasmin_aus-station_r1i1p1_CSIRO-MnCh-wrt-1986-2005-Scl_v1_mon_seasavg-clim.csv",
            "hurs15_aus-station_r1i1p1_CSIRO-MnCh-wrt-1986-2005-Scl_v1_mon_seasavg-clim.csv",
            "pan-evap_aus-station_r1i1p1_CSIRO-MnCh-wrt-1986-2005-Scl_v1_mon_seassum-clim.csv"
        ]
        
        for file in climate_files:
            file_path = os.path.join(self.data_dir, file)
            if not os.path.exists(file_path):
                print(f"⚠️ Skipping {file} - file not found")
                continue
                
            print(f"  Processing {file}...")
            df = pd.read_csv(file_path)
            
            # Determine climate type from filename
            if "tasmax" in file:
                climate_type = "max_temperature"
                unit = "°C"
            elif "tasmin" in file:
                climate_type = "min_temperature"
                unit = "°C"
            elif "tas" in file and "tasmax" not in file and "tasmin" not in file:
                climate_type = "avg_temperature"
                unit = "°C"
            elif "hurs" in file:
                climate_type = "humidity"
                unit = "%"
            elif "pan-evap" in file:
                climate_type = "evaporation"
                unit = "mm"
            else:
                climate_type = "climate"
                unit = "unknown"
            
            unified_data = []
            for _, row in df.iterrows():
                lat_val = row['LAT']
                lon_val = row['LON']
                if pd.isna(lat_val) or pd.isna(lon_val):
                    continue
                    
                metadata = {
                    "station_name": row.get('STATION_NAME', ''),
                    "station_id": row.get('STN_ID', ''),
                    "model": row.get('MODEL', ''),
                    "rcp": row.get('RCP', ''),
                    "annual": row.get('Annual', 0),
                    "djf": row.get('DJF', 0),
                    "mam": row.get('MAM', 0),
                    "jja": row.get('JJA', 0),
                    "son": row.get('SON', 0),
                    "climate_type": climate_type,
                    "unit": unit
                }
                
                unified_data.append({
                    "dataset_type": "climate",
                    "name": f"{climate_type}_{row.get('STATION_NAME', 'Unknown')}",
                    "lat": row['LAT'],
                    "lon": row['LON'],
                    "value": row.get('Annual', 0),
                    "unit": unit,
                    "metadata": json.dumps(metadata)
                })
            
            self.save_to_sql(unified_data, f"climate_{climate_type}.sql")
    
    def process_wind_data(self):
        """Process wind observation data"""
        print("💨 Processing wind data...")
        input_file = os.path.join(self.data_dir, "wind-observations.csv")
        
        if not os.path.exists(input_file):
            print(f"⚠️ Skipping wind data - file not found: {input_file}")
            return
            
        df = pd.read_csv(input_file)
        
        # Filter out invalid coordinates
        df = df[(df['latitude'] != 0.0) & (df['longitude'] != 0.0)]
        
        unified_data = []
        for _, row in df.iterrows():
            # Skip records with missing coordinates
            if pd.isna(row['latitude']) or pd.isna(row['longitude']):
                continue
                
            metadata = {
                "location_description": row.get('location_description', ''),
                "gust_speed": row.get('gust_speed', 0),
                "wind_direction": row.get('wind_direction', 0),
                "wind_direction_cardinal": row.get('wind_direction_cardinal', ''),
                "date_time": str(row.get('date_time', ''))
            }
            
            unified_data.append({
                "dataset_type": "wind",
                "name": f"wind_{row.get('location_description', 'Unknown')}",
                "lat": row['latitude'],
                "lon": row['longitude'],
                "value": row.get('average_wind_speed', 0),
                "unit": "m/s",
                "metadata": json.dumps(metadata)
            })
        
        self.save_to_sql(unified_data, "wind_observations.sql")
    
    def process_vegetation_data(self):
        """Process vegetation zone data"""
        print("🌿 Processing vegetation data...")
        input_file = os.path.join(self.data_dir, "VegetationZones_718376949849166399.csv")
        
        if not os.path.exists(input_file):
            print(f"⚠️ Skipping vegetation data - file not found: {input_file}")
            return
            
        df = pd.read_csv(input_file)
        
        unified_data = []
        for _, row in df.iterrows():
            metadata = {
                "zone": row.get('Zone', ''),
                "type": row.get('Type', ''),
                "shape_area": row.get('SHAPE_area', 0),
                "shape_length": row.get('SHAPE_len', 0),
                "link": row.get('Link', '')
            }
            
            # For vegetation, we'll use the zone center (approximate)
            # Since this is polygon data, we'll need to calculate centroids
            # For now, using placeholder coordinates
            unified_data.append({
                "dataset_type": "vegetation",
                "name": f"vegetation_{row.get('Type', 'Unknown')}",
                "lat": -37.8136,  # Melbourne approximate center
                "lon": 144.9631,
                "value": row.get('SHAPE_area', 0),
                "unit": "m²",
                "metadata": json.dumps(metadata)
            })
        
        self.save_to_sql(unified_data, "vegetation_zones.sql")
    
    def save_to_sql(self, data, filename):
        """Save unified data to SQL file"""
        if not data:
            print(f"  ⚠️ No data to save for {filename}")
            return
            
        output_file = os.path.join(self.output_dir, filename)
        
        sql_statements = []
        for i, item in enumerate(data):
            sql = f"""INSERT INTO datasets (dataset_type, name, lat, lon, value, unit, metadata) 
                     VALUES ('{item['dataset_type']}', '{item['name'].replace("'", "''")}', 
                             {item['lat']}, {item['lon']}, {item['value']}, '{item['unit']}', 
                             '{item['metadata'].replace("'", "''")}');"""
            sql_statements.append(sql)
        
        with open(output_file, "w", encoding="utf-8") as f:
            f.write("\n".join(sql_statements))
        
        print(f"  ✅ Saved {len(data)} records to {filename}")

if __name__ == "__main__":
    processor = DatasetProcessor()
    processor.process_all_datasets() 