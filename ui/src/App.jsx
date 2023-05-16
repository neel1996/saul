import { useEffect, useState } from "react";
import { RouterProvider } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

import { router } from "./Routes";
import FullPageLoader from "./components/FullPageLoader";
import { initializeFirebase } from "./firebase";
import { useAxiosInterceptor } from "./useAxiosInterceptor";

function App() {
    const [showLoader, setShowLoader] = useState(false);

    useEffect(() => {
        initializeFirebase();
    }, []);

    useAxiosInterceptor({ setShowLoader });

    return (
        <>
            <ToastContainer limit={3} />
            <RouterProvider router={router}></RouterProvider>
            {showLoader && <FullPageLoader />}
        </>
    );
}

export default App;
