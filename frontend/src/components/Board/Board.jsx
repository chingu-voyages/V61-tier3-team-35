import Tile from "./Tile"

export default function Board({ className, board }) {

    return (
        <div className={`flex flex-col gap-2 ${className}`}>
            {board.map((row, rowIndex) => (
                <div key={rowIndex} className="grid grid-cols-5 gap-2 w-fit">
                    {row.map((tile, colIndex) => (
                        <Tile key={colIndex} letter={tile.letter} status={tile.status}/>
                    ))}
                </div>
            ))}

        </div>
    )
}