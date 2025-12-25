import React from "react";

export function highlightText(text: string, query: string) {
  const q = query.trim();
  if (!q) return text;

  const lower = text.toLowerCase();
  const needle = q.toLowerCase();

  const nodes: React.ReactNode[] = [];
  let i = 0;

  while (true) {
    const idx = lower.indexOf(needle, i);
    if (idx === -1) break;

    if (idx > i) nodes.push(text.slice(i, idx));

    nodes.push(
      <span
        key={`${idx}-${needle}`}
        className="rounded-[3px] bg-[#0b4ea2] px-0.5 text-white"
      >
        {text.slice(idx, idx + needle.length)}
      </span>
    );

    i = idx + needle.length;
  }

  if (i < text.length) nodes.push(text.slice(i));
  return nodes;
}
