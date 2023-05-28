import BackgroundWrapper from "@components/BackgroundWrapper";
import { faFilePdf } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { uploadDocument } from "@services/fileUpload";
import React, { useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { useNavigate } from "react-router-dom";

export default function DocumentUpload() {
    const navigate = useNavigate();

    const onDrop = useCallback((acceptedFiles) => {
        const file = acceptedFiles[0];
        uploadDocument(file)
            .then((data) => {
                const { checksum } = data;
                navigate(`/document-qa?checksum=${checksum}`);
            })
            .catch((error) => {
                console.log(error);
            });
    }, [navigate]);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
    });

    return (
        <BackgroundWrapper>
            <div
                {...getRootProps()}
                className="flex mx-auto my-auto select-none"
            >
                <input {...getInputProps()} />
                {isDragActive ? (
                    <div>Drop the files here ...</div>
                ) : (
                    <div className="text-center mx-auto my-auto p-14 bg-gray-700 rounded-lg shadow-sm text-base font-semibold text-gray-400 select-none border-2 border-dashed border-slate-600">
                        <FontAwesomeIcon
                            icon={faFilePdf}
                            size="6x"
                            className="my-6 text-indigo-500"
                        />
                        <div className="my-4">
                            Drag and drop some files here, or click to select
                            files
                        </div>
                    </div>
                )}
            </div>
        </BackgroundWrapper>
    );
}
