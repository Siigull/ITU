import { useState, useEffect } from "react";

// hook to get window size dynamically
const use_window_size = () => {
    // Initialize state with undefined width/height so server and client renders match
    // Learn more here: https://joshwcomeau.com/react/the-perils-of-rehydration/
    const [windowSize, setWindowSize] = useState({
        width: undefined,
        height: undefined
    });

    useEffect(() => {
        // only execute all the code below in client side
        if (typeof window !== "undefined") {
            // Handler to call on window resize
            function handleResize() {
                // Set window width/height to state
                setWindowSize({
                    width: window.innerWidth,
                    height: window.innerHeight
                });
            }

            // Add event listener
            window.addEventListener("resize", handleResize);

            // Call handler right away so state gets updated with initial window size
            handleResize();

            // Remove event listener on cleanup
            return () => window.removeEventListener("resize", handleResize);
        }
    }, []); // Empty array ensures that effect is only run on mount
    return windowSize;
};

export default use_window_size;

export const generate_user_auth = (user_id, game_id, token) => {
    var data = {}
    data.user_id = user_id
    data.game_id = game_id
    data.token = token 

    return data
}

export const game_type_to_string = (game) => {
    let ret_string = game.game_type;
    if(ret_string == "" || ret_string == undefined) {
        return ""
    }

    if(game.trump_suit == '0') {
        ret_string += " acorn";

    } else if(game.trump_suit == '1') {
        ret_string += " bells";

    } else if(game.trump_suit == '2') {
        ret_string += " hearts";

    } else {
        ret_string += " leaves";
    }

    return ret_string
}

export const butt_style = "bg-white hover:bg-gray-100 text-gray-800 font-semibold py-2 px-4 border border-gray-400 rounded shadow"

export const get_card_src = (card_val) => {
    return `/cards/${card_val.suit}-${card_val.spec}.png`
}