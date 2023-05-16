import Login from "@pages/Login";
import {
    Route,
    createBrowserRouter,
    createRoutesFromElements,
} from "react-router-dom";

export const router = createBrowserRouter(
    createRoutesFromElements(
        <Route path="/">
            <Route path="/login" element={<Login />} />
        </Route>
    )
);
