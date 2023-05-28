import { post } from "@root/AxiosInstance";

export const uploadDocument = async (file) => {
    const formData = new FormData();
    formData.append("file", file);

    return await post("/upload", formData)
        .then((data) => {
            return data;
        })
        .catch((error) => {
            throw new Error(error);
        });
};
