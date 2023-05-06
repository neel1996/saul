import { RouterProvider } from "react-router-dom";
import { router } from "./Routes";

function App() {
    return (
        <div className="bg-slate-800">
            <RouterProvider router={router}></RouterProvider>
        </div>
    );
}

export default App;
