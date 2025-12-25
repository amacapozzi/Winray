import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

export function PaletteHeader({
  onClose,
}: {
  title: string;
  onClose: () => void;
}) {
  return (
    <div className="flex justify-end px-4">
      <Button
        variant="ghost"
        size="icon"
        className="h-8 w-8 text-white/70 hover:bg-white/5 hover:text-white/90"
        type="button"
        onClick={onClose}
      >
        <X className="h-4.5 w-4.5" />
      </Button>
    </div>
  );
}
