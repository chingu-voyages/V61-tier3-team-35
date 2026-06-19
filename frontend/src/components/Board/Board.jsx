import Tile from "./Tile"

export default function Board({ className }) {
    const rows = Array(6).fill(null);
    const cols = Array(5).fill(null);

    return (
        <div className={`flex flex-col gap-2 ${className}`}>
            {rows.map((row, rowIndex) => (
                <div key={rowIndex} className="grid grid-cols-5 gap-2 w-fit">
                    {cols.map((col, colIndex) => (
                        <Tile key={colIndex} />
                    ))}
                </div>
            ))}

        </div>
    )
}