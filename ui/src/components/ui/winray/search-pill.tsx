import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search, RotateCcw, ChevronUp, ChevronDown } from "lucide-react";

export function SearchPill({
  value,
  onChange,
  positionLabel,
  onReset,
  onPrev,
  onNext,
}: {
  value: string;
  onChange: (v: string) => void;
  positionLabel: string;
  onReset?: () => void;
  onPrev?: () => void;
  onNext?: () => void;
}) {
  return (
    <div className="px-5 pt-4 pb-3">
      <div className="flex items-center gap-3 rounded-lg border border-white/10 bg-[#2a2b2d] px-3 py-2">
        <div className="text-white/55">
          <Search className="h-4.5 w-4.5" />
        </div>

        <Input
          id="winray-search"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="Search files…"
          className="h-9 flex-1 border-0 bg-transparent  text-[16px] text-white/90 placeholder:text-white/35 focus-visible:ring-0 focus-visible:ring-offset-0"
        />

        <div className="text-[13px] text-white/40">{positionLabel}</div>

        <div className="flex items-center overflow-hidden rounded-md border border-white/10 bg-[#232426]">
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-10 rounded-none text-white/55 hover:bg-white/5 hover:text-white/80"
            type="button"
            onClick={onReset}
          >
            <RotateCcw className="h-4 w-4" />
          </Button>

          <div className="h-9 w-px bg-white/10" />

          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-10 rounded-none text-white/55 hover:bg-white/5 hover:text-white/80"
            type="button"
            onClick={onPrev}
          >
            <ChevronUp className="h-4 w-4" />
          </Button>

          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-10 rounded-none text-white/55 hover:bg-white/5 hover:text-white/80"
            type="button"
            onClick={onNext}
          >
            <ChevronDown className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
