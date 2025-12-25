import type { FileResult } from "./types";
import { SearchPill } from "./search-pill";
import { ResultList } from "./result-list";
import { usePaletteController } from "./use-palette-controller";
import { useEffect, useState } from "react";

declare global {
  interface Window {
    goSearch?: (q: string) => void;
    goOpen?: (path: string, kind: string) => void;
    goHide?: () => void;
    goStartIndexing?: () => void;

    // Go -> UI
    setResults?: (results: FileResult[]) => void;
    appendResults?: (results: FileResult[]) => void;
    setLoading?: (loading: boolean) => void;

    // Go puede llamarte para focus/clear
    focusSearch?: () => void;
    clearSearch?: () => void;
  }
}

export function WinrayPaleette() {
  const [results, setResults] = useState<FileResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // Bridge: Go -> React
  useEffect(() => {
    window.setResults = (r) => {
      setResults(r);
    };
    window.appendResults = (r) => {
      setResults((prev) => {
        const existingIds = new Set(prev.map((item) => item.id));
        const newResults = r.filter((item) => !existingIds.has(item.id));
        return [...prev, ...newResults];
      });
    };
    window.setLoading = (loading) => setIsLoading(loading);
    window.focusSearch = () => {
      document.querySelector<HTMLInputElement>("#winray-search")?.focus();
    };
    window.clearSearch = () => {};

    const init = () => {
      setTimeout(() => {
        window.goStartIndexing?.();
      }, 150);
    };

    if (document.readyState === "complete") {
      init();
    } else {
      window.addEventListener("load", init);
    }
  }, []);

  const onClose = () => window.goHide?.();

  const onOpen = (item: FileResult) => {
    window.goOpen?.(item.path, item.kind || "File");
  };

  const controller = usePaletteController({
    results,
    onOpen,
    onClose,
  });

  useEffect(() => {
    window.goSearch?.(controller.query);
  }, [controller.query]);

  const goPrev = () =>
    controller.setActiveIndex(Math.max(controller.activeIndex - 1, 0));
  const goNext = () =>
    controller.setActiveIndex(
      Math.min(
        controller.activeIndex + 1,
        Math.max(controller.filtered.length - 1, 0)
      )
    );

  const onReset = () => controller.setQuery("");

  const positionLabel = controller.positionLabel;

  return (
    <div className="bg-[#0e0f10] fixed inset-0 flex items-center justify-center">
      <div
        className={[
          "w-[620px] rounded-none",
          "bg-[#1f2022]/95 backdrop-blur-xl",
          "shadow-[0_20px_60px_rgba(0,0,0,.55)]",
          "border border-white/10",
        ].join(" ")}
      >
        <SearchPill
          value={controller.query}
          onChange={controller.setQuery}
          positionLabel={positionLabel}
          onReset={onReset}
          onPrev={goPrev}
          onNext={goNext}
        />

        <style>{`#winray-search{}`}</style>

        {isLoading && results.length === 0 ? (
          <div className="px-5 py-20 flex flex-col items-center justify-center">
            <div className="w-8 h-8 border-2 border-white/20 border-t-white/60 rounded-full animate-spin mb-3" />
            <div className="text-sm text-white/40">Indexing files...</div>
          </div>
        ) : (
          <ResultList
            items={controller.filtered}
            query={controller.query}
            activeIndex={controller.activeIndex}
            onActiveIndexChange={controller.setActiveIndex}
            onOpen={onOpen}
          />
        )}
        {isLoading && results.length > 0 && (
          <div className="px-5 py-3 flex items-center gap-2 text-sm text-white/40">
            <div className="w-4 h-4 border-2 border-white/20 border-t-white/60 rounded-full animate-spin" />
            <span>Indexing...</span>
          </div>
        )}
      </div>
    </div>
  );
}
