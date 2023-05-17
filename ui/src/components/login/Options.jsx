import { faGithub, faGoogle } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { loginWithGoogle, loginWithGithub } from "@services/login";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

import FullPageLoader from "../FullPageLoader";

export default function Options() {
    const [showLoader, setShowLoader] = useState(false);
    const navigate = useNavigate();

    const loginHandler = (loginPromise) => {
        setShowLoader(true);

        loginPromise
            .then(() => {
                setShowLoader(false);
                return navigate("/dashboard");
            })
            .catch((err) => {
                setShowLoader(false);
                throw err;
            });
    };

    const loginOptions = [
        {
            name: "google",
            icon: faGoogle,
            login: async () => {
                loginHandler(loginWithGoogle());
            },
        },
        {
            name: "github",
            icon: faGithub,
            login: async () => {
                loginHandler(loginWithGithub());
            },
        },
    ];

    return (
        <>
            {showLoader && <FullPageLoader showLoader={showLoader} />}
            <div data-testid="login-options" className="my-4">
                <div className="text-center text-base font-semibold text-gray-500 select-none">
                    Sign in with
                </div>
                <div
                    data-testid="login-options-cta"
                    className="flex items-center align-middle justify-center gap-8 my-4"
                >
                    {loginOptions.map((option) => (
                        <div
                            key={option.name}
                            data-testid={`login-with-${option.name}`}
                            className="my-auto p-3 text-base border border-gray-600 transition-all rounded-full shadow-lg bg-gray-600 cursor-pointer hover:shadow-sm text-indigo-400"
                            onClick={option.login}
                        >
                            <FontAwesomeIcon icon={option.icon} size="2x" />
                        </div>
                    ))}
                </div>
            </div>
        </>
    );
}
