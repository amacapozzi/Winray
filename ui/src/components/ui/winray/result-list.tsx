import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import type { FileResult } from "./types";
import { ResultRow } from "./result-row";

export function ResultList({
  items,
  query,
  activeIndex,
  onActiveIndexChange,
  onOpen,
}: {
  items: FileResult[];
  query: string;
  activeIndex: number;
  onActiveIndexChange: (i: number) => void;
  onOpen: (item: FileResult) => void;
}) {
  return (
    <div className="px-0 pt-3">
      <ScrollArea className="h-80">
        <div className="pb-2">
          {items.length === 0 ? (
            <div className="px-5 py-10 text-center text-sm text-white/40">
              No results
            </div>
          ) : (
            items.map((it, i) => (
              <div key={it.id}>
                <ResultRow
                  item={it}
                  query={query}
                  active={i === activeIndex}
                  onHover={() => onActiveIndexChange(i)}
                  onOpen={() => onOpen(it)}
                />
                <Separator className="bg-white/10" />
              </div>
            ))
          )}
        </div>
      </ScrollArea>
    </div>
  );
}
