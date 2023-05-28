import { io } from "socket.io-client";

export const socket = io(
    process.env.VITE_SOCKET_BASE_URL || "http://localhost:8080/socket.io/",
    {
        extraHeaders: {
            Authorization: `Bearer ${sessionStorage.getItem("authToken")}`,
        },
        cors: {
            origin: "*",
        },
    }
);
