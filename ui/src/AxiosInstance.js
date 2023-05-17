import axios from "axios";

export const AxiosInstance = axios.create({
    baseURL: process.env.VITE_BASE_URL || "http://localhost:8080/api/saul/v1",
    timeout: 10000,
    cancelToken: axios.CancelToken.source().token,
});

export const get = async (url, config, cancelToken) => {
    return await AxiosInstance.get(url, {
        ...config,
        headers: {
            ...config?.headers,
            Authorization: "Bearer " + sessionStorage.getItem("authToken"),
        },
        cancelToken,
    })
        .then((response) => {
            return response.data;
        })
        .catch((error) => {
            throw new Error(error);
        });
};

export const post = async (url, data, config, cancelToken) => {
    return await AxiosInstance.post(url, data, {
        ...config,
        headers: {
            ...config?.headers,
            Authorization: "Bearer " + sessionStorage.getItem("authToken"),
        },
        cancelToken,
    })
        .then((response) => {
            return response.data;
        })
        .catch((error) => {
            throw new Error(error);
        });
};
