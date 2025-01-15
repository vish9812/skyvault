import Thumbnail from "@/components/Thumbnail";
import FormattedDateTime from "@/components/FormattedDateTime";
import { Link } from "react-router";
import { FileModel } from "@/lib/models";
import ActionDropdown from "./ActionDropdown";
import utils from "@/lib/utils";

interface Props {
  file: FileModel;
}

const Card = ({ file }: Props) => {
  return (
    <Link to={file.url} target="_blank" className="file-card">
      <div className="flex justify-between">
        <Thumbnail
          type={file.type}
          extension={file.extension}
          url={file.url}
          className="!size-20"
          imageClassName="!size-11"
        />

        <div className="flex flex-col items-end justify-between">
          <ActionDropdown file={file} />
          <p className="body-1">{utils.prettySize(file.sizeBytes)}</p>
        </div>
      </div>

      <div className="file-card-details">
        <p className="subtitle-2 line-clamp-1">{file.name}</p>
        <FormattedDateTime
          date={file.createdAt}
          className="body-2 text-light-100"
        />
        <p className="caption line-clamp-1 text-light-200">
          By: {file.ownerId}
        </p>
      </div>
    </Link>
  );
};
export default Card;
