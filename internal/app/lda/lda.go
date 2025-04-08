package lda

import (
	//"errors"
	"fmt"
	"math"

	"github.com/lda_api/internal/app/entity"
	"github.com/lda_api/internal/app/repository"
	"github.com/sirupsen/logrus"
	//"gonum.org/v1/gonum/mat"
)

// TODO -> разобраться со структурой:
//
//	-> переписать методы, в которых использовалась ConversionData под ConvertDataContext
type LDA struct {
	db *repository.DataBase
	//ConvAllData entity.ConversionData
	Alpha         float64   `json:"alpha"`
	Beta          float64   `json:"beta"`
	AccuracyModel float64   `json:"accuracy"`
	ShiftingModel float64   `json:"shifting"`
	X             []float64 `json:"x"`
	Y             []float64 `json:"y"`
}

func New(db *repository.DataBase) *LDA {
	return &LDA{
		Alpha: 0.0,
		Beta:  0.0,
		db:    db,
	}
}

func (lda *LDA) ConvertData(dataUsers entity.ContextData) entity.ConvertDataContext {
	var dataContext entity.ConvertDataContext
	for i := 0; i < len(dataUsers.Data); i++ {
		diffCoef := (float64)(dataUsers.Data[i].LoanAmount) / float64(dataUsers.Data[i].IncomeAnnum)
		diffCoef = (float64(dataUsers.Data[i].IncomeAnnum) / diffCoef) / 1000
		if dataUsers.Data[i].SelfEmployed == " Yes" {
			diffCoef = diffCoef * 0.25
		} else {
			diffCoef = diffCoef * 0.1
		}
		var dataConv entity.ConversionData
		dataConv.ImportancecCoefficient = diffCoef
		dataConv.Rating = float64(dataUsers.Data[i].CibilScore)
		if dataUsers.Data[i].LoanStatus == "Approved" {
			dataConv.Class = 1
		} else {
			dataConv.Class = 0
		}
		dataContext.Data = append(dataContext.Data, dataConv)

	}
	return dataContext
}

func GetArray(data entity.ConvertDataContext) ([]float64, []float64, []float64) {
	var dataRating, dataCoef, dataClass []float64

	for i := 0; i < len(data.Data); i++ {
		dataRating = append(dataRating, data.Data[i].Rating)
		dataCoef = append(dataCoef, data.Data[i].ImportancecCoefficient)
		dataClass = append(dataClass, float64(data.Data[i].Class))
	}
	return dataRating, dataCoef, dataClass
}
func GetMean(data []float64) float64 {
	if data == nil {
		return 0.0
	}
	var sum float64 = 0
	for i := 0; i < len(data); i++ {
		sum = sum + data[i]
	}

	return sum / float64(len(data))
}

func GetDispersion(data []float64, mean float64) float64 {
	if data == nil {
		return 0.0
	}
	var sum float64 = 0.0
	for i := 0; i < len(data); i++ {
		var temp float64 = data[i] - mean
		sum = sum + math.Pow(temp, 2)
	}
	return sum / (float64(len(data) - 1))
}

func GetCovariation(dataX []float64, meanX float64, dataY []float64, meanY float64) float64 {
	if dataX == nil && dataY == nil {
		logrus.Fatal("Data is Empty")
		return 0.0
	}
	var sum float64 = 0.0
	for i := 0; i < len(dataX); i++ {
		sum = sum + (dataX[i]-meanX)*(dataY[i]-meanY)
	}
	return sum / (float64(len(dataY) - 1))
}

func GetCovariationMatrix(mean []float64, dataCoef, dataRating []float64) [][]float64 {
	var covariationMatrix [][]float64
	covariation := GetCovariation(dataCoef, mean[1], dataRating, mean[0])
	for i := 0; i < len(mean); i++ {
		var temp []float64
		for j := 0; j < len(mean); j++ {
			if i == j {
				if i == 0 {
					temp = append(temp, GetDispersion(dataRating, mean[i]))
				} else {
					temp = append(temp, GetDispersion(dataCoef, mean[i]))
				}
			} else {
				temp = append(temp, covariation)
			}
		}
		covariationMatrix = append(covariationMatrix, temp)
	}
	return covariationMatrix
}

func PlusMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	for i := 0; i < len(leftOperand); i++ {
		for j := 0; j < len(leftOperand[i]); j++ {
			leftOperand[i][j] = leftOperand[i][j] + rightOperand[i][j]
		}
	}
	return leftOperand
}

func SubtractionMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	for i := 0; i < len(leftOperand); i++ {
		for j := 0; j < len(leftOperand[i]); j++ {
			leftOperand[i][j] = leftOperand[i][j] - rightOperand[i][j]
		}
	}
	return leftOperand
}
func ProdMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	var resultMatrix [][]float64
	var sum float64 = 0.0

	for i := 0; i < len(leftOperand); i++ {
		var temp []float64
		for j := 0; j < len(rightOperand[i]); j++ {
			sum = 0
			for k := 0; k < len(leftOperand); k++ {
				sum = sum + (leftOperand[i][k] * rightOperand[k][j])
			}
			temp = append(temp, sum)
		}
		resultMatrix = append(resultMatrix, temp)
	}
	return resultMatrix
}

func GetDet(matrix [][]float64) float64 {
	var prod float64 = 1
	var prod2 float64 = 1
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if i == j {
				prod = prod * matrix[i][j]
			} else {
				prod2 = prod2 * matrix[i][j]
			}
		}
	}
	return (prod - prod2)
}

func CheckUnitMatrix(inverseMatrix, W [][]float64) bool {
	result := ProdMatrix(W, inverseMatrix)
	var count int = 0
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i]); j++ {
			if result[i][j] >= 0.9 && result[i][j] <= 1 && i == j {
				count++
			} else {
				if int(result[i][j]) == 0 && i != j {
					count++
				}
			}
		}
	}

	if count == (len(result) * len(result[0])) {
		return true
	} else {
		return false
	}
}

func InverseMatrix(W [][]float64) [][]float64 {
	var det float64 = GetDet(W)

	var updW [][]float64
	for i := 1; i <= len(W); i += 1 {
		var tempForAdd []float64
		for j := 1; j <= len(W[0]); j += 1 {
			temp := math.Pow((-1), float64(i+j)) * W[len(W)-i][len(W[0])-j]
			tempForAdd = append(tempForAdd, temp)
		}
		updW = append(updW, tempForAdd)
	}

	var result [][]float64
	for i := 0; i < len(W); i++ {
		var tempForAdd []float64
		for j := 0; j < len(W[i]); j++ {
			temp := (1 / det) * updW[i][j]
			tempForAdd = append(tempForAdd, temp)
		}
		result = append(result, tempForAdd)
	}

	if CheckUnitMatrix(result, W) {
		return result
	} else {
		logrus.Fatal("Матрица не является обратной")
		return nil
	}

}

func GetW(covariationMatrixFirstClass, covariationMatrixSecondClass [][]float64) [][]float64 {
	if covariationMatrixFirstClass == nil && covariationMatrixSecondClass == nil {
		return nil
	}
	W := PlusMatrix(covariationMatrixFirstClass, covariationMatrixSecondClass)
	for i := 0; i < len(covariationMatrixFirstClass); i++ {
		for j := 0; j < len(covariationMatrixFirstClass[i]); j++ {
			W[i][j] = W[i][j] / 2
		}
	}
	return W
}

func GetB(covariationMatrixAllData, W [][]float64) [][]float64 {
	return SubtractionMatrix(covariationMatrixAllData, W)
}

func GetS(W, B [][]float64) [][]float64 {
	W = InverseMatrix(W)
	return ProdMatrix(W, B)
}

func FindScale(a, b float64, dataCoef, dataRating []float64) ([]float64, float64, float64, float64) {

	var score []float64 = nil

	for score == nil { //|| GetDispersion(score, GetMean(score)) > 1.0{
		for i := 0; i < len(dataCoef); i++ {
			temp := a*dataCoef[i] + b*dataRating[i]
			score = append(score, temp)
		}
		disp := GetDispersion(score, GetMean(score))
		if disp > 1.1 {
			a = a / math.Sqrt(disp)
			b = b / math.Sqrt(disp)
			score = nil
		} else {
			break
		}
	}
	return score, GetMean(score), a, b

}

func (lda *LDA) PredictModel(dataUser entity.UserData) (string, float64) {
	var context entity.ContextData
	context.Data = append(context.Data, dataUser)
	dataConvertion := lda.ConvertData(context)
	score := lda.Alpha*dataConvertion.Data[0].ImportancecCoefficient + lda.Beta*dataConvertion.Data[0].Rating + lda.ShiftingModel
	var predict string
	if score > 0 {
		logrus.Info("Дискриминационная оценка данного пользователя: ", score)
		logrus.Info("Кредит будет успешно одобрен")
		predict = "Approved"
	} else {
		logrus.Info("Дискриминационная оценка данного пользователя: ", score)
		logrus.Info("Кредит будет отклонен")
		predict = "Rejected"
	}
	return predict, score
}

func (lda *LDA) GetAccuracy(dataCoef, dataRating, classData []float64) {
	countTryAnswer := 0.0
	for i := 0; i < len(dataCoef); i++ {
		score := lda.Alpha*dataCoef[i] + lda.Beta*dataRating[i] + lda.ShiftingModel
		var predict int
		if score > 0 {
			predict = 1
		} else {
			predict = 0
		}
		if predict == int(classData[i]) {
			countTryAnswer += 1
		}
	}
	lda.AccuracyModel = (countTryAnswer / float64(len(classData))) * 100.0
}

func GetCombination(vectorsForCombination [][]float64) [][]float64 {
	var resultComb [][]float64

	for i := 0; i < len(vectorsForCombination); i++ {
		for j := 0; j < len(vectorsForCombination[i]); j++ {
			tempA := vectorsForCombination[i][j]
			for k := 0; k < len(vectorsForCombination); k++ {
				var comb []float64
				for l := 0; l < len(vectorsForCombination[k]); l++ {
					if k != i || l != j {
						tempB := vectorsForCombination[k][l]
						comb = append(comb, tempA)
						comb = append(comb, tempB)
					}
				}
				resultComb = append(resultComb, comb)
			}
		}
	}
	return resultComb
}
func (lda *LDA) FindProjection(meanClass []float64) float64 {
	return lda.Alpha*meanClass[0] + lda.Beta*meanClass[1]
}
func (lda *LDA) FitModel() error {
	// инициализация массивов средних значений для классов и всего набора данных
	var meanValueFirstClass []float64
	var meanValueSecondClass []float64
	var meanValueAllData []float64
	/// инициализация ковариационных матриц для классов и всего набора данных
	//var covariationMatrixAllData [][]float64
	var covariationMatrixFirstClass [][]float64
	var covariationMatrixSecondClass [][]float64
	// инициализация ковариационных матриц S, W, B
	//var S []float64
	var W [][]float64
	//var B [][]float64

	// Получение данных для получения значения точности обученной модели
	dataTestUsers, err := lda.db.Data().SelectTestData()
	if err != nil {
		return err
	}
	convDataTest := lda.ConvertData(dataTestUsers)
	testDataRating, testDataImportancecCoefficient, testDataClass := GetArray(convDataTest)
	// Получение всех необходимых значений для всего набора данных
	dataAllUsers, err := lda.db.Data().SelectAllLearnData()
	if err != nil {
		return err
	}
	convAllData := lda.ConvertData(dataAllUsers)
	allDataRating, allDataImportancecCoefficient, _ := GetArray(convAllData)
	meanValueAllData = append(meanValueAllData, GetMean(allDataRating))
	meanValueAllData = append(meanValueAllData, GetMean(allDataImportancecCoefficient))

	/*covariationMatrixAllData = GetCovariationMatrix(
		meanValueAllData,
		allDataImportancecCoefficient,
		allDataRating,
	)
	*/
	// Получение данных для класса "Approved"
	firstClass, err := lda.db.Data().SelectingDataByClass("Approved")
	if err != nil {
		return err
	}
	dataFirstClass := lda.ConvertData(firstClass)
	firstDataRating, firstDataImportancecCoefficient, _ := GetArray(dataFirstClass)
	meanValueFirstClass = append(meanValueFirstClass, GetMean(firstDataRating))
	meanValueFirstClass = append(meanValueFirstClass, GetMean(firstDataImportancecCoefficient))

	covariationMatrixFirstClass = GetCovariationMatrix(
		meanValueFirstClass,
		firstDataImportancecCoefficient,
		firstDataRating,
	)

	// Получение данных для класса "Rejected"
	secondClass, err := lda.db.Data().SelectingDataByClass("Rejected")
	if err != nil {
		return err
	}
	dataSecondClass := lda.ConvertData(secondClass)
	secondDataRating, secondDataImportancecCoefficient, _ := GetArray(dataSecondClass)
	meanValueSecondClass = append(meanValueSecondClass, GetMean(secondDataRating))
	meanValueSecondClass = append(meanValueSecondClass, GetMean(secondDataImportancecCoefficient))

	covariationMatrixSecondClass = GetCovariationMatrix(
		meanValueSecondClass,
		secondDataImportancecCoefficient,
		secondDataRating,
	)

	// Произведение расчетов для ковариационных матриц S, B, W

	W = GetW(covariationMatrixFirstClass, covariationMatrixSecondClass)

	//B = GetB(covariationMatrixAllData, W)
	/*temp := GetS(W, B)
	//конвертация данных под формат метода получения собственных значений и векторов
	for i := 0; i < len(temp); i++ {
		for j := 0; j < len(temp[i]); j++ {
			S = append(S, temp[i][j])
		}
	}
	//инициализация объекта для получения собственных значений и векторов
	A := mat.NewDense(2, 2, S)
	var eig mat.Eigen
	if !eig.Factorize(A, mat.EigenRight) {
		err := errors.New("error in calculating eigenvalues")
		return err
	}

	// Получаем собственные значения
	values := eig.Values(nil)
	if values == nil {
		err := errors.New("couldn't get eigenvalues")
		return err
	}

	// Получаем собственные векторы
	var vectors mat.CDense
	eig.VectorsTo(&vectors)
	var vectorsForCombination [][]float64
	for i := 0; i < 2; i++ {
		var temp []float64
		for j := 0; j < 2; j++ {
			temp = append(temp, float64(real(vectors.At(i, j))))
		}
		vectorsForCombination = append(vectorsForCombination, temp)
	}
	combination := GetCombination(vectorsForCombination)

	for i := 0; i < len(combination); i++ {

		tempAccuracy := 0.0
		lda.Score,
			lda.ConstDiscrimitation,
			lda.Alpha,
			lda.Beta = FindScale(combination[i][0], combination[i][1], allDataImportancecCoefficient, allDataRating)
		if lda.Alpha != 0.0 && lda.Beta != 0.0 {
			tempAccuracy = lda.GetAccuracy(testDataImportancecCoefficient, testDataRating, testDataClass)
		}
		if lda.AccuracyModel < tempAccuracy {
			lda.AccuracyModel = tempAccuracy
		}
	}*/

	W = InverseMatrix(W)
	var differenceMean [][]float64

	for i := 0; i < len(meanValueFirstClass); i++ {
		tempValue := meanValueFirstClass[i] - meanValueSecondClass[i]
		var tmp []float64
		tmp = append(tmp, tempValue)
		differenceMean = append(differenceMean, tmp)
	}
	omega := ProdMatrix(W, differenceMean)
	fmt.Println(omega)
	lda.Alpha = omega[0][0]
	lda.Beta = omega[1][0]
	lda.ShiftingModel = ((lda.FindProjection(meanValueFirstClass) + lda.FindProjection(meanValueSecondClass)) / 2.0) * -1
	lda.GetAccuracy(testDataRating, testDataImportancecCoefficient, testDataClass)
	logrus.Info("The model has been successfully trained.")
	logrus.Info("Model accuracy: ", lda.AccuracyModel)
	//logrus.Infof("Коэффициенты модели %f и %f", lda.Alpha, lda.Beta)

	return nil
}
