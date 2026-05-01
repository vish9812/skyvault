/**
 * Format bytes to human-readable format
 * @param bytes - The number of bytes
 * @returns A formatted string representing the size
 */
function size(bytes?: number) {
  if (bytes === undefined) return "-";
  if (bytes === 0) return "0 B";

  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
}

/**
 * Format a date string to a more readable format
 * @param dateString - The date string to format
 * @returns A formatted date string
 */
function date(dateString?: string) {
  if (!dateString) return "-";
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

/**
 * Get the initials from a name
 * @param name - The name to extract initials from
 * @returns A string of initials
 */
function initials(name: string) {
  const parts = name.split(" ");
  const firstChar = String.fromCodePoint(parts[0].codePointAt(0)!);

  if (parts.length < 2) return firstChar.toUpperCase();

  const lastChar = String.fromCodePoint(
    parts[parts.length - 1].codePointAt(0)!,
  );

  return `${firstChar}${lastChar}`.toUpperCase();
}

const Format = {
  initials,
  size,
  date,
} as const;

export default Format;
