import Numeric;
import Data.Char;

opcodes xs = let
				toUppercase = map toUpper
				toOpcodeStr = (++) "0x"
			 in
				map (toOpcodeStr . toUppercase . flip showHex "") xs

toInstrTest f xs = do
					timings <- readFile f >>= return . lines
					ops <- return $ opcodes xs
					zipped <- return $ zip ops timings
					return $ map (\(opcode, timing) -> "RunInstrAndAssertTimings(" ++ opcode ++ ", " ++ timing ++ ", nil, t)")  zipped
