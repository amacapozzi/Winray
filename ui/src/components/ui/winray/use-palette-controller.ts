import { useEffect, useMemo, useState } from "react";
import type { FileResult } from "./types";

type Options = {
  results: FileResult[];
  onOpen: (item: FileResult) => void;
  onClose: () => void;
};

export function usePaletteController({ results, onOpen, onClose }: Options) {
  const [query, setQuery] = useState("");
  const [activeIndex, setActiveIndex] = useState(0);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return results.slice(0, 60);

    const out: FileResult[] = [];
    for (const r of results) {
      const hay = `${r.name} ${r.path}`.toLowerCase();
      if (hay.includes(q)) out.push(r);
      if (out.length >= 60) break;
    }
    return out;
  }, [query, results]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    if (activeIndex >= filtered.length) setActiveIndex(0);
  }, [filtered.length, activeIndex]);

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") {
        e.preventDefault();
        onClose();
        return;
      }
      if (e.key === "ArrowDown") {
        e.preventDefault();
        setActiveIndex((v) =>
          Math.min(v + 1, Math.max(filtered.length - 1, 0))
        );
        return;
      }
      if (e.key === "ArrowUp") {
        e.preventDefault();
        setActiveIndex((v) => Math.max(v - 1, 0));
        return;
      }
      if (e.key === "Enter") {
        e.preventDefault();
        const item = filtered[activeIndex];
        if (item) onOpen(item);
      }
    }

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [activeIndex, filtered, onClose, onOpen]);

  const positionLabel = filtered.length
    ? `${activeIndex + 1}/${filtered.length}`
    : "0/0";

  return {
    query,
    setQuery,
    filtered,
    activeIndex,
    setActiveIndex,
    positionLabel,
  };
}
