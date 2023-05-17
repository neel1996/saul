import {
    getAuth,
    signInWithPopup,
    GoogleAuthProvider,
    GithubAuthProvider,
} from "firebase/auth";

import { loginWithGoogle, loginWithGithub } from "./firebaseAuth";
import { login } from "./userLogin";

jest.mock("./userLogin");
jest.mock("@root/AxiosInstance");
jest.mock("firebase/auth");

describe("firebaseAuth", () => {
    afterAll(() => {
        jest.clearAllMocks();
        jest.resetAllMocks();
    });

    beforeEach(() => {
        signInWithPopup.mockResolvedValue({
            user: {
                uid: "123",
                displayName: "John Doe",
                email: "test@test.com",
                photoURL: "https://test.com",
                accessToken: "testAuthToken",
            },
        });
    });

    it("should login with Google", async () => {
        await loginWithGoogle();

        expect(getAuth).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledWith(
            getAuth(),
            new GoogleAuthProvider()
        );
        expect(login).toHaveBeenCalledTimes(1);
        expect(login).toHaveBeenCalledWith({
            userId: "123",
            name: "John Doe",
            email: "test@test.com",
            avatar: "https://test.com",
        });
    });

    it("should login with Github", async () => {
        await loginWithGithub();

        expect(getAuth).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledWith(
            getAuth(),
            new GithubAuthProvider()
        );
        expect(login).toHaveBeenCalledTimes(1);
        expect(login).toHaveBeenCalledWith({
            userId: "123",
            name: "John Doe",
            email: "test@test.com",
            avatar: "https://test.com",
        });
    });
});
