import axios from "axios";

export const AxiosInstance = axios.create({
    baseURL: process.env.VITE_BASE_URL || "http://localhost:8080/api/saul/v1",
    timeout: 10000,
    headers: { Authorization: "Bearer " + sessionStorage.getItem("authToken") },
    cancelToken: axios.CancelToken.source().token,
});
