import BackgroundImage from "@assets/background-overlay.png";
import { PropTypes } from "prop-types";
import React from "react";

export default function BackgroundWrapper({ children }) {
    return (
        <div
            className="w-full h-full flex mx-auto my-auto bg-indigo-500"
            style={{
                backgroundImage: `url(${BackgroundImage})`,
                backgroundSize: "70%",
                backgroundRepeat: "no-repeat",
                backgroundPosition: "center",
            }}
        >
            {children}
        </div>
    );
}

BackgroundWrapper.propTypes = {
    children: PropTypes.node.isRequired,
};
