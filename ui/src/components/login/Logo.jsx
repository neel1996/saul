import { faBolt } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";

export default function Logo() {
    return (
        <div data-testid="logo">
            <div className="mx-auto flex align-middle items-center justify-center gap-4 select-none">
                <div className="p-4 rounded-lg shadow-md bg-gray-600 text-indigo-400">
                    <FontAwesomeIcon icon={faBolt} size="2x" />
                </div>
                <div className="block items-center">
                    <div className="font-bold text-gray-300 text-6xl">Saul</div>
                    <div className="font-sans text-sm font-semibold text-indigo-300">
                        Know what is in your document
                    </div>
                </div>
            </div>
            <div className="border-b border-gray-600 w-11/12 mx-auto my-6"></div>
        </div>
    );
}
