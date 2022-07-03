import Image from "next/image";
import { useEffect, useRef, useState } from "react";
import useResizeObserver from "use-resize-observer";
import FileInput from "./FileInput";

const Canvas = ({ onImageChange, image, className }) => {
  const ref = useRef();
  const [canvasSize, setCanvasSize] = useState({ width: 1, height: 1 });
  const { width, height } = useResizeObserver({ ref });

  useEffect(() => {
    if (!!ref.current) {
      const rect = ref.current.children[0].getBoundingClientRect();
      setCanvasSize({ width: rect.width, height: rect.height });
    }
  }, [width, height, ref]);

  return (
    <div className={"flex flex-col gap-4 items-center p-3 " + className}>
      {image ? (
        <div className="relative">
          <div ref={ref}>
            <Image
              src={image.src}
              width={image.width}
              height={image.height}
              alt="Image being sonified"
              className="peer absolute top-0 left-0 z-0 rounded-md shadow-md"
            />
          </div>
          <canvas
            id="pixelsound"
            width={canvasSize.width}
            height={canvasSize.height}
            className="absolute top-0 left-0 z-10 rounded-md"
          ></canvas>
        </div>
      ) : (
        <div className="animate-pulse bg-slate-100 w-48 h-48 rounded-md shadow-md"></div>
      )}
      <FileInput onChange={onImageChange} accept=".jpg,.jpeg,.png,.gif">
        Select an image
      </FileInput>
    </div>
  );
};

export default Canvas;
