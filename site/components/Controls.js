import FileInput from "./FileInput";
import Modes from "./modes/Modes";

const Controls = ({ onImageChange, onAudioChange, onModeChange }) => {
  return (
    <div className="flex flex-col max-w-lg items-center mx-auto">
      <div className="flex flex-wrap truncate gap-4 p-3">
        <FileInput onChange={onImageChange} accept=".jpg,.jpeg,.png">
          Select an image
        </FileInput>
        <FileInput onChange={onAudioChange} accept=".mp3">
          Select an audio file
        </FileInput>
      </div>
      {/* <Modes onChange={onModeChange}></Modes> */}
    </div>
  );
};

export default Controls;
