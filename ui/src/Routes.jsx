import {
    Route,
    createBrowserRouter,
    createRoutesFromElements,
} from "react-router-dom";
import Login from "./pages/Login";

export const router = createBrowserRouter(
    createRoutesFromElements(<Route path="/login" element={<Login />} />)
);
