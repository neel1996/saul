import { loginWithGoogle, loginWithGithub } from "./firebaseAuth";
import {
    getAuth,
    signInWithPopup,
    GoogleAuthProvider,
    GithubAuthProvider,
} from "firebase/auth";

jest.mock("firebase/auth");

describe("firebaseAuth", () => {
    it("should login with Google", () => {
        loginWithGoogle();

        expect(getAuth).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledWith(
            getAuth(),
            new GoogleAuthProvider()
        );
    });

    it("should login with Github", () => {
        loginWithGithub();

        expect(getAuth).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledTimes(1);
        expect(signInWithPopup).toHaveBeenCalledWith(
            getAuth(),
            new GithubAuthProvider()
        );
    });
});
