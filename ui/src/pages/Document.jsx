import {
    faPaperPlane,
    faRobot,
    faUserAlt,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useRef, useState } from "react";
import { ThreeDots } from "react-loader-spinner";
import { useSearchParams } from "react-router-dom";
import io from "socket.io-client";

export default function Document() {
    const [searchParams] = useSearchParams();
    const checksum = searchParams.get("checksum");

    const [message, setMessage] = useState("");
    const [conversations, setConversations] = useState([]);
    const [loading, setLoading] = useState(false);

    let socket = useRef(null);
    useEffect(() => {
        socket.current = io(
            process.env.VITE_SOCKET_BASE_URL || "http://localhost:8080",
            {
                path: "/saul/socket.io",
            }
        );

        socket.current.on("connect", () => {
            socket.current.emit("join", checksum);
        });

        socket.current.on("disconnect", () => {
            console.log("Disconnected from server");
        });
    }, [checksum]);

    useEffect(() => {
        socket.current.on("answer", (data) => {
            if (!data) return;

            const parsedData = JSON.parse(data);
            setLoading(false);

            setConversations((prev) => [
                ...prev,
                {
                    message: parsedData.answer,
                    sender: "bot",
                },
            ]);
        });
    }, []);

    const sendMessage = () => {
        if (message === "") return;

        setLoading(true);
        setConversations((prev) => [
            ...prev,
            {
                message,
                sender: "user",
            },
        ]);

        socket.current.emit(
            "message",
            JSON.stringify({
                documentId: checksum,
                question: message,
            })
        );

        setMessage("");
    };

    return (
        <div className="mx-auto my-auto select-none p-6 bg-gray-700 w-full h-full">
            <div className="h-1/8">
                <div className="text-center select-none text-white font-sans font-semibold text-xl">
                    Post your queries and get answers!
                </div>
                <div className="border-b border-gray-800 w-full block my-3"></div>
            </div>
            <div id="chats" className="h-5/6 overflow-auto">
                {conversations.map((conversation, index) => {
                    return (
                        <div
                            key={index}
                            className={`flex ${
                                conversation.sender === "user"
                                    ? "justify-end"
                                    : "justify-start"
                            }`}
                        >
                            <FontAwesomeIcon
                                icon={
                                    conversation.sender === "user"
                                        ? faUserAlt
                                        : faRobot
                                }
                                className={`${
                                    conversation.sender === "user"
                                        ? "text-blue-500 bg-gray-800"
                                        : "text-gray-800 bg-blue-400"
                                } my-auto mx-2 p-2  rounded-full`}
                            />
                            <div
                                className={`${
                                    conversation.sender === "user"
                                        ? "bg-blue-500"
                                        : "bg-gray-500"
                                } rounded-md px-3 py-2 my-2 mx-2 text-white`}
                            >
                                {conversation.message}
                            </div>
                        </div>
                    );
                })}
                {loading ? (
                    <div className="flex justify-start">
                        <div className="bg-gray-500 rounded-md px-3 py-2 my-2 mx-2 text-white">
                            <ThreeDots
                                height="30"
                                width="30"
                                radius="9"
                                color="#ffffff"
                                ariaLabel="three-dots-loading"
                                wrapperStyle={{}}
                                wrapperClassName=""
                                visible={true}
                            />
                        </div>
                    </div>
                ) : null}
            </div>
            <form
                className="flex w-full bg-gray-500 rounded-md h-20"
                onSubmit={(e) => {
                    e.preventDefault();

                    sendMessage();
                }}
            >
                <input
                    type="text"
                    className="w-full bg-gray-500 rounded-md outline-none text-gray-200 font-sans text-lg px-2"
                    placeholder="Type your message here..."
                    value={message}
                    onChange={(e) => {
                        setMessage(e.target.value);
                    }}
                />
                <input type="submit" id="send" hidden></input>
                <label
                    htmlFor="send"
                    className="w-24 bg-gray-800 flex text-white justify-center items-center rounded-r-md shadow-md hover:bg-gray-700 cursor-pointer"
                >
                    <FontAwesomeIcon
                        icon={faPaperPlane}
                        className="w-1/3 h-1/3 mx-auto my-auto"
                    />
                </label>
            </form>
        </div>
    );
}
