import Board from "../components/Board/Board"
import Keyboard from "../components/Keyboard/Keyboard"


export default function(){
    return (
        <div className="flex flex-col items-center">
            <Board />
            <Keyboard />
        </div>
    )
}