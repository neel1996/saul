import { AxiosInstance } from "@root/AxiosInstance";

export const validateUser = async (request) => {
    const { data } = await AxiosInstance.post("/user/validate", request);

    return data;
};
