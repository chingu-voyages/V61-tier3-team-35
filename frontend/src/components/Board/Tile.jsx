const statuses =
{
    correct: "text-white border-none bg-lime-500",
    present: "text-white border-none bg-yellow-400",
    absent: "text-white border-none bg-gray-400",
}


export default function Tile({ letter, status }) {


    return (
        <div className={`uppercase  flex items-center justify-center text-2xl font-bold md:size-14 sm:size-12 size-11 border  ${statuses[status]} bg-gray-50/50 ${letter ? "border-black border-2" : "border-gray-400"}`}>{letter}</div>
    )
}