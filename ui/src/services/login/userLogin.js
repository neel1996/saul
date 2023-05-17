import { post } from "@root/AxiosInstance";

export const login = async (request) => {
    return await post("/login", request)
        .then((data) => {
            if (data.authToken) {
                sessionStorage.setItem("authToken", data.authToken);
                return data;
            }

            throw new Error("Invalid credentials");
        })
        .catch((error) => {
            throw new Error(error);
        });
};
