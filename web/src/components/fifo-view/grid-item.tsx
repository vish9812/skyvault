interface GridItemProps {
  type: "file" | "folder";
  name: string;
  preview?: string;
}

function GridItem(props: GridItemProps) {
  const getFileIcon = () => {
    if (props.type === "folder") {
      return "folder";
    }

    // Determine file type icon based on name/extension
    if (props.name) {
      const extension = props.name.split(".").pop()?.toLowerCase();

      if (
        ["jpg", "jpeg", "png", "gif", "svg", "webp"].includes(extension || "")
      ) {
        return "image";
      } else if (["mp4", "webm", "avi", "mov"].includes(extension || "")) {
        return "movie";
      } else if (["mp3", "wav", "ogg"].includes(extension || "")) {
        return "audio_file";
      } else if (["doc", "docx", "pdf", "txt"].includes(extension || "")) {
        return "description";
      }
    }

    return "insert_drive_file";
  };

  return (
    <div class="bg-white rounded-lg border border-border shadow-sm hover:border-primary hover:shadow-md transition-all group">
      {/* File/folder icon or preview */}
      <div class="flex-center h-28 border-b border-border bg-bg-subtle relative">
        {props.preview ? (
          <img
            src={`data:image/png;base64,${props.preview}`}
            alt={props.name}
            class="object-cover h-full w-full"
          />
        ) : (
          <span
            class={`material-symbols-outlined text-5xl ${
              props.type === "folder" ? "text-primary" : "text-neutral-light"
            }`}
          >
            {getFileIcon()}
          </span>
        )}
        <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <span class="material-symbols-outlined text-neutral-light cursor-pointer hover:text-primary bg-white rounded-full shadow-sm p-1">
            more_vert
          </span>
        </div>
      </div>

      {/* File/folder info */}
      <div class="p-3">
        <div class="text-sm font-medium text-neutral truncate">
          {props.name}
        </div>
      </div>
    </div>
  );
}

export default GridItem;
