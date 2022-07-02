import Image from "next/image";
import FileInput from "./FileInput";

const Canvas = ({ imgSrc = undefined, className }) => {
  return (
    <div className={"flex flex-col gap-4 items-center py-3 " + className}>
      {imgSrc ? (
        <div>
          <Image src={imgSrc} alt="Image being sonified"></Image>
          <canvas id="pixelsound" className="rounded-md shadow-md"></canvas>
        </div>
      ) : (
        <div className="animate-pulse bg-slate-100 w-48 h-48 rounded-md shadow-md"></div>
      )}
      <FileInput>Select an image</FileInput>
    </div>
  );
};

export default Canvas;
