import { RouterProvider } from "react-router-dom";
import { router } from "./Routes";
import { initializeFirebase } from "./firebase";

function App() {
    initializeFirebase();

    return <RouterProvider router={router}></RouterProvider>;
}

export default App;
