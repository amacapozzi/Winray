import { Badge } from "@/components/ui/badge";
import type { FileResult } from "./types";
import { highlightText } from "./highlight";

function formatTimeAgo(timestamp?: number): string {
  if (!timestamp) return "";

  const now = Date.now() / 1000;
  const diff = now - timestamp;

  if (diff < 60) return "hace un momento";
  if (diff < 3600) {
    const minutes = Math.floor(diff / 60);
    return `hace ${minutes} min${minutes > 1 ? "s" : ""}`;
  }
  if (diff < 86400) {
    const hours = Math.floor(diff / 3600);
    return `hace ${hours} hora${hours > 1 ? "s" : ""}`;
  }
  if (diff < 604800) {
    const days = Math.floor(diff / 86400);
    return `hace ${days} día${days > 1 ? "s" : ""}`;
  }
  if (diff < 2592000) {
    const weeks = Math.floor(diff / 604800);
    return `hace ${weeks} semana${weeks > 1 ? "s" : ""}`;
  }
  if (diff < 31536000) {
    const months = Math.floor(diff / 2592000);
    return `hace ${months} mes${months > 1 ? "es" : ""}`;
  }
  const years = Math.floor(diff / 31536000);
  return `hace ${years} año${years > 1 ? "s" : ""}`;
}

export function ResultRow({
  item,
  query,
  active,
  onHover,
  onOpen,
}: {
  item: FileResult;
  query: string;
  active: boolean;
  onHover: () => void;
  onOpen: () => void;
}) {
  const timeAgo = item.lastAccessTime
    ? formatTimeAgo(item.lastAccessTime)
    : null;

  return (
    <button
      type="button"
      onMouseEnter={onHover}
      onClick={onOpen}
      className={[
        "w-full text-left px-5 py-3 transition-colors",
        active ? "bg-white/5" : "bg-transparent hover:bg-white/5",
      ].join(" ")}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="text-[16px] leading-[1.15] text-white/90">
            {highlightText(item.name, query)}
          </div>
          {item.metaLeft && (
            <div className="mt-1 truncate text-[13px] text-white/45">
              {highlightText(item.metaLeft, query)}
            </div>
          )}
          {timeAgo && (
            <div className="mt-1 text-[13px] text-white/45">{timeAgo}</div>
          )}
        </div>

        {item.kind ? (
          <Badge
            variant="secondary"
            className="shrink-0 bg-white/5 text-white/55 border border-white/10 rounded "
          >
            {item.kind}
          </Badge>
        ) : null}
      </div>
    </button>
  );
}
