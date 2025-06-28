"use client";
import { useState, useEffect } from "react";
import styles from "./DatasetSelector.module.scss";

interface DatasetInfo {
  type: string;
  name: string;
  description: string;
  count: number;
  date_range?: string;
  value_range?: string;
}

interface DatasetSelectorProps {
  onDatasetChange: (datasetType: string) => void;
  selectedDataset: string;
}

export default function DatasetSelector({ onDatasetChange, selectedDataset }: DatasetSelectorProps) {
  const [datasets, setDatasets] = useState<DatasetInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchDatasetTypes();
  }, []);

  const fetchDatasetTypes = async () => {
    try {
      const BASE_URL = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") || "http://localhost:8080";
      const response = await fetch(`${BASE_URL}/datasets/types`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      setDatasets(data);
    } catch (error) {
      console.error("❌ Error fetching dataset types:", error);
      setError(error instanceof Error ? error.message : "Failed to fetch datasets");
    } finally {
      setLoading(false);
    }
  };

  const handleDatasetChange = (datasetType: string) => {
    onDatasetChange(datasetType);
  };

  if (loading) {
    return <div className={styles.loading}>Loading datasets...</div>;
  }

  if (error) {
    return <div className={styles.error}>Error: {error}</div>;
  }

  return (
    <div className={styles.selector}>
      <h3 className={styles.title}>🌍 Select Dataset</h3>
      <div className={styles.grid}>
        {datasets.map((dataset) => (
          <div
            key={dataset.type}
            className={`${styles.card} ${selectedDataset === dataset.type ? styles.selected : ''}`}
            onClick={() => handleDatasetChange(dataset.type)}
          >
            <div className={styles.icon}>
              {getDatasetIcon(dataset.type)}
            </div>
            <div className={styles.content}>
              <h4 className={styles.name}>{dataset.name}</h4>
              <p className={styles.description}>{dataset.description}</p>
              <div className={styles.stats}>
                <span className={styles.count}>{dataset.count.toLocaleString()} records</span>
                {dataset.value_range && (
                  <span className={styles.range}>Range: {dataset.value_range}</span>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function getDatasetIcon(datasetType: string): string {
  switch (datasetType) {
    case "meteorite":
      return "☄️";
    case "climate":
      return "🌡️";
    case "wind":
      return "💨";
    case "vegetation":
      return "🌿";
    case "infrastructure":
      return "🏗️";
    case "fire":
      return "🔥";
    default:
      return "📍";
  }
}