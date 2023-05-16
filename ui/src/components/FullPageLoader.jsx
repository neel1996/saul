import PropTypes from "prop-types";
import React from "react";
import { Dna } from "react-loader-spinner";

export default function FullPageLoader({ showLoader }) {
    return (
        showLoader && (
            <div className="fixed top-0 left-0 w-full h-full z-50 flex justify-center items-center bg-gray-900 bg-opacity-80">
                <Dna
                    visible={showLoader}
                    height="150"
                    width="150"
                    ariaLabel="dna-loading"
                    wrapperClass="dna-wrapper"
                />
            </div>
        )
    );
}

FullPageLoader.propTypes = {
    showLoader: PropTypes.bool.isRequired,
};
