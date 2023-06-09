import BackgroundImage from "@assets/background-overlay.png";
import React from "react";

import Logo from "./Logo";
import Options from "./Options";

export default function Login() {
    return (
        <div
            data-testid="login"
            className="flex mx-auto my-auto w-full h-full bg-indigo-500"
        >
            <div
                className="flex mx-auto my-auto w-full h-full"
                style={{
                    backgroundImage: `url(${BackgroundImage})`,
                    backgroundSize: "70%",
                    backgroundRepeat: "no-repeat",
                    backgroundPosition: "center",
                }}
            >
                <div className="bg-gray-700 p-6 rounded-lg shadow-lg w-11/12 md:w-1/2 xl:w-1/3 sm:w-11/12 mx-auto my-auto border border-gray-600">
                    <Logo />
                    <Options />
                </div>
            </div>
        </div>
    );
}
