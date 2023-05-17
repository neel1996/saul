import { toast } from "react-toastify";

import { AxiosInstance } from "./AxiosInstance";

export const useAxiosInterceptor = ({ setShowLoader }) => {
    AxiosInstance.interceptors.request.use(
        (config) => {
            setShowLoader(true);
            return config;
        },
        (error) => {
            setShowLoader(false);
            return Promise.reject(error);
        }
    );

    AxiosInstance.interceptors.response.use(
        (response) => {
            setShowLoader(false);
            return response;
        },
        (error) => {
            setShowLoader(false);
            toast(error?.response?.data?.message || "Something went wrong!", {
                toastId: "errorToast",
                position: "bottom-center",
                autoClose: 5000,
                hideProgressBar: false,
                closeOnClick: true,
                pauseOnHover: false,
                draggable: true,
                progress: 100,
                theme: "dark",
                progressStyle: {
                    background: "#fe5454",
                },
            });
            return Promise.reject(error);
        }
    );
};
