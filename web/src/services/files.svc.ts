import consts from "@/lib/consts";

const filesURLPvt = consts.configs.baseAPIPvt + "/files";

async function uploadFile(file: File) {
  // Return fake promise
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({ name: file.name });
    }, 2000);
  });

  // const formData = new FormData();
  // formData.append("file", file);

  // const res = await fetch(filesURLPvt, {
  //   method: "POST",
  //   headers: consts.headers.auth(),
  //   body: formData,
  // });

  // if (res.ok) {
  //   return res.json();
  // } else {
  //   const errorText = await res.text();
  //   throw new Error(errorText);
  // }
}

const filesSvc = {
  uploadFile,
};

export default filesSvc;
