"use client";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import DatasetList from "../../../components/DatasetList";
import SearchForm from "../../../components/SearchForm";
import styles from "./page.module.scss";

interface DatasetInfo {
  type: string;
  name: string;
  description: string;
  count: number;
  date_range?: string;
  value_range?: string;
}

export default function DatasetPage() {
  const params = useParams();
  const router = useRouter();
  const datasetType = params.type as string;
  
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [datasetInfo, setDatasetInfo] = useState<DatasetInfo | null>(null);

  useEffect(() => {
    fetchDatasetInfo();
  }, [datasetType]);

  const fetchDatasetInfo = async () => {
    try {
      const BASE_URL = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") || "http://localhost:8080";
      const response = await fetch(`${BASE_URL}/datasets/types`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const datasets = await response.json();
      const currentDataset = datasets.find((d: DatasetInfo) => d.type === datasetType);
      
      if (!currentDataset) {
        throw new Error(`Dataset type '${datasetType}' not found`);
      }
      
      setDatasetInfo(currentDataset);
    } catch (error) {
      console.error("❌ Error fetching dataset info:", error);
      setError(error instanceof Error ? error.message : "Failed to fetch dataset info");
    }
  };

  const searchDatasets = async (query: Record<string, string | number>) => {
    setLoading(true);
    setError(null);
    setResults([]); 

    try {
      const BASE_URL = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") || "http://localhost:8080";
      const url = new URL(`${BASE_URL}/datasets`);

      // Add dataset type
      url.searchParams.append("type", datasetType);

      Object.entries(query).forEach(([key, value]) => {
        if (value !== "" && value !== undefined) {
          // Map old parameter names to new ones
          const newKey = key === "mass_min" ? "value_min" : 
                        key === "mass_max" ? "value_max" : key;
          url.searchParams.append(newKey, String(value).trim());
        }
      });

      console.log("🔍 Query Parameters:", query);
      console.log("🌍 Fetching from:", url.toString());

      const response = await fetch(url.toString());

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      console.log("✅ API Response:", data);

      // Check if the response is an error object
      if (data.error) {
        throw new Error(data.error);
      }

      if (!Array.isArray(data)) {
        console.error("❌ Unexpected response format:", data);
        throw new Error(
          `Unexpected API response format: Expected an array, got ${typeof data}. Response: ${JSON.stringify(data)}`
        );
      }

      setResults(data);
    } catch (error) {
      console.error("❌ Error fetching datasets:", error);
      setError(
        error instanceof Error ? error.message : "An unexpected error occurred."
      );
    } finally {
      setLoading(false);
    }
  };

  const handleBackToHome = () => {
    router.push("/");
  };

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <div className={styles.error}>
          <h2>Error: {error}</h2>
          <button onClick={handleBackToHome} className={styles.backButton}>
            ← Back to Home
          </button>
        </div>
      </div>
    );
  }

  if (!datasetInfo) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.loading}>Loading dataset information...</div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <button onClick={handleBackToHome} className={styles.backButton}>
          ← Back to Categories
        </button>
        <div className={styles.titleSection}>
          <h1 className={styles.title}>
            {getDatasetIcon(datasetType)} {datasetInfo.name}
          </h1>
          <p className={styles.description}>{datasetInfo.description}</p>
          <div className={styles.stats}>
            <span className={styles.count}>{datasetInfo.count.toLocaleString()} records</span>
            {datasetInfo.value_range && (
              <span className={styles.range}>Value range: {datasetInfo.value_range}</span>
            )}
          </div>
        </div>
      </div>

      <div className={styles.content}>
        <SearchForm onSearch={searchDatasets} loading={loading} datasetType={datasetType} />
        {error && <p className={styles.errorMessage}>{error}</p>}
        <DatasetList results={results} datasetType={datasetType} />
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