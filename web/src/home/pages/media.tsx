import filesSvc from "@/services/files.svc";
import React, { useEffect } from "react";
import { useParams } from "react-router";
import { FileModel } from "@/lib/models";
import Card from "@/components/Card";

const Media = () => {
  const { mediaType } = useParams();
  // const searchText = ((await searchParams)?.query as string) || "";
  // const sort = ((await searchParams)?.sort as string) || "";

  // const types = getFileTypesParams(type) as FileType[];

  const [files, setFiles] = React.useState<FileModel[]>([]);

  useEffect(() => {
    async function getFiles() {
      try {
        setFiles([]);
        const f = await filesSvc.getFiles(null);
        if (!ignore) {
          setFiles(f);
        }
      } catch (err) {
        console.error(err);
      }
    }

    let ignore = false;
    getFiles();

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <div className="page-container">
      <section className="w-full">
        <h1 className="h1 capitalize">{mediaType}</h1>

        <div className="total-size-section">
          <p className="body-1">
            Total: <span className="h5">0 MB</span>
          </p>

          <div className="sort-container">
            <p className="body-1 hidden text-light-200 sm:block">Sort by:</p>

            {/* <Sort /> */}
          </div>
        </div>
      </section>

      {files.length > 0 ? (
        <section className="file-list">
          {files.map((f: FileModel) => (
            <Card key={f.id} file={f} />
          ))}
        </section>
      ) : (
        <p className="empty-list">No files uploaded</p>
      )}
    </div>
  );
};

export default Media;
