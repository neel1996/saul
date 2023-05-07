import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGithub, faGoogle } from "@fortawesome/free-brands-svg-icons";

export default function Options() {
    const loginOptions = [
        {
            name: "Google",
            icon: faGoogle,
        },
        {
            name: "github",
            icon: faGithub,
        },
    ];

    return (
        <div data-testid="login-options" className="my-4">
            <div className="text-center text-base font-semibold text-gray-500">
                Sign in with
            </div>
            <div
                data-testid="login-options-cta"
                className="flex items-center align-middle justify-center gap-8 my-4"
            >
                {loginOptions.map((option) => (
                    <div
                        key={option.name}
                        className="my-auto p-3 text-base border border-gray-600 transition-all rounded-full shadow-lg bg-gray-600 cursor-pointer hover:shadow-sm text-indigo-400"
                    >
                        <FontAwesomeIcon icon={option.icon} size="2x" />
                    </div>
                ))}
            </div>
        </div>
    );
}