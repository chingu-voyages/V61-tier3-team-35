export default function SelectMode() {
    return (
        <article className="flex gap-2">
            <button className="text-xs flex justify-center gap-2 sm:w-25 w-fit text-white bg-taupe-500 rounded-md p-2"> <span className="sm:block hidden">Daily</span>  📅</button>
            <button className="text-xs flex justify-center gap-2 sm:w-25 w-fit text-white bg-taupe-500 rounded-md p-2"> <span className="sm:block hidden">Unlimited</span>  ♾️</button>
        </article>
    )
}