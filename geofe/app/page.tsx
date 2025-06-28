"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import styles from "./page.module.scss";
import Link from "next/link";

interface DatasetInfo {
  type: string;
  name: string;
  description: string;
  count: number;
  date_range?: string;
  value_range?: string;
}

const DATASET_CARDS = [
  {
    type: "meteorite",
    icon: "☄️",
    title: "Meteorites",
    desc: "Global meteorite impact and discovery data",
    color: "#7f5af0"
  },
  {
    type: "climate",
    icon: "🌡️",
    title: "Climate Stations",
    desc: "Australian climate station data with temperature, humidity, and evaporation",
    color: "#ff6f3c"
  },
  {
    type: "wind",
    icon: "💨",
    title: "Wind Observations",
    desc: "Wind speed and direction observations",
    color: "#43d9ad"
  },
  {
    type: "vegetation",
    icon: "🌿",
    title: "Vegetation Zones",
    desc: "Vegetation classification and area data",
    color: "#38b000"
  },
  {
    type: "infrastructure",
    icon: "🏗️",
    title: "Infrastructure",
    desc: "Infrastructure and utility data",
    color: "#f15bb5"
  },
  {
    type: "fire",
    icon: "🔥",
    title: "Fire Projections",
    desc: "Fire risk and projection data",
    color: "#ff3c6f"
  }
];

export default function HomePage() {
  const [datasets, setDatasets] = useState<DatasetInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const [stats, setStats] = useState<any>(null);

  useEffect(() => {
    fetchDatasetTypes();
    fetchStats();
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

  const fetchStats = async () => {
    try {
      const BASE_URL = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") || "http://localhost:8080";
      const res = await fetch(`${BASE_URL}/datasets/types`);
      const data = await res.json();
      setStats(data);
    } catch (e) {
      setStats(null);
    }
  };

  const totalRecords = stats ? stats.reduce((sum: number, d: any) => sum + d.count, 0) : 0;
  const datasetCount = stats ? stats.length : 0;
  const latestType = stats && stats[0] ? stats[0].name : "-";

  const handleDatasetSelect = (datasetType: string) => {
    router.push(`/dataset/${datasetType}`);
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.loading}>Loading datasets...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <div className={styles.error}>Error: {error}</div>
      </div>
    );
  }

  return (
    <div className={styles.darkBg}>
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <h1 className={styles.title}>GeoGO: Explore Global Geospatial Data</h1>
          <p className={styles.subtitle}>
            Meteorites, Climate, Wind, Vegetation, Infrastructure, Fire & more
          </p>
          <Link href="#datasets" className={styles.exploreBtn}>
            Explore Datasets
          </Link>
        </div>
        <div className={styles.heroGlobe}>
          {/* can be animated*/}
          <img src="/globe.svg" alt="Globe" />
        </div>
      </section>

      <section className={styles.statsRow}>
        <div className={styles.statCard}>
          <span className={styles.statLabel}>Total Datasets</span>
          <span className={styles.statValue}>{datasetCount}</span>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statLabel}>Total Records</span>
          <span className={styles.statValue}>{totalRecords.toLocaleString()}</span>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statLabel}>Latest Dataset</span>
          <span className={styles.statValue}>{latestType}</span>
        </div>
      </section>

      <section className={styles.datasetsSection} id="datasets">
        <div className={styles.datasetsGrid}>
          {DATASET_CARDS.map(card => {
            const stat = stats && stats.find((d: any) => d.type === card.type);
            return (
              <Link
                href={`/dataset/${card.type}`}
                key={card.type}
                className={styles.datasetCard}
                style={{
                  background: `linear-gradient(135deg, ${card.color}33 0%, #18181b 100%)`,
                  borderColor: card.color
                }}
              >
                <div className={styles.cardIcon}>{card.icon}</div>
                <div className={styles.cardTitle}>{card.title}</div>
                <div className={styles.cardDesc}>{card.desc}</div>
                {stat && (
                  <div className={styles.cardStat}>{stat.count.toLocaleString()} records</div>
                )}
              </Link>
            );
          })}
        </div>
      </section>

      <section className={styles.mapRow}>
        <div className={styles.mapPreview}>
          {/* Placeholder for mini map preview, can be replaced with Leaflet or SVG */}
          <img src="/globe.svg" alt="Mini Map" />
        </div>
        <div className={styles.didYouKnow}>
          <h3>Did you know?</h3>
          <p>Australia is home to the world's largest meteorite impact structure, the Vredefort crater!</p>
        </div>
      </section>

      <footer className={styles.footer}>
        <span>GeoGO &copy; {new Date().getFullYear()} &mdash; All data open & free</span>
        <span><a href="https://github.com/klea/GeoGO" target="_blank" rel="noopener noreferrer">GitHub</a></span>
      </footer>
    </div>
  );
}
