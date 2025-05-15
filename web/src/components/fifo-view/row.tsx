interface RowProps {
  type: "file" | "folder";
  name: string;
  size?: number;
  updatedAt?: string;
}

function Row(props: RowProps) {
  return (
    <div class="flex items-center py-3 px-4 hover:bg-gray-50 group relative">
      <span class="material-symbols-outlined text-primary mr-3 text-2xl">
        folder
      </span>
      <div class="flex-1 min-w-0">
        <div class="text-base font-semibold text-gray-900 truncate">
          {props.name}
        </div>
      </div>
    </div>
  );
}

export default Row;
