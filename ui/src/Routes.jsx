import { Login, DocumentUpload, Document } from "@pages";
import {
    Route,
    createBrowserRouter,
    createRoutesFromElements,
    Navigate,
} from "react-router-dom";

export const router = createBrowserRouter(
    createRoutesFromElements(
        <Route path="/">
            <Route path="/login" element={<Login />} />
            <Route path="/document-upload" element={<DocumentUpload />} />
            <Route path="/document-qa" element={<Document />} />
            <Route path="*" element={<Navigate to="/login" replace />} />
        </Route>
    )
);
