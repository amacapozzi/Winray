export type FileResult = {
  id: string;
  name: string;
  path: string;
  kind?: "App" | "Folder" | "File" | "Shortcut";
  metaLeft?: string;
  metaRight?: string;
  lastAccessTime?: number;
};
